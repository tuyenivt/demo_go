package auction

import (
	"context"
	"sync"
	"time"

	"real-time-bidding/cache"
	"real-time-bidding/models"
)

// AuctionService handles the RTB auction process
type AuctionService struct {
	cache *cache.RedisCache
}

// NewAuctionService creates a new auction service instance
func NewAuctionService(cache *cache.RedisCache) *AuctionService {
	return &AuctionService{
		cache: cache,
	}
}

// ProcessBidRequest handles an incoming bid request
func (s *AuctionService) ProcessBidRequest(ctx context.Context, req *models.BidRequest) (*models.BidResponse, error) {
	startTime := time.Now()

	// Get user segments
	_, err := s.cache.GetUserSegments(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// Get eligible campaigns (simplified for example)
	eligibleCampaigns := []string{"campaign1", "campaign2", "campaign3"}

	// Process bids concurrently
	bidChan := make(chan *models.BidResponse, len(eligibleCampaigns))
	var wg sync.WaitGroup

	for _, campaignID := range eligibleCampaigns {
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if bid, err := s.evaluateBid(ctx, req, cid); err == nil {
				bidChan <- bid
			}
		}(campaignID)
	}

	// Wait for all bids to be processed
	go func() {
		wg.Wait()
		close(bidChan)
	}()

	// Find the highest bid
	var bestBid *models.BidResponse
	for bid := range bidChan {
		if bestBid == nil || bid.Price > bestBid.Price {
			bestBid = bid
		}
	}

	if bestBid == nil {
		return nil, nil // No eligible bids
	}

	// Update campaign budget
	if err := s.cache.UpdateCampaignBudget(ctx, bestBid.CampaignID, bestBid.Price); err != nil {
		return nil, err
	}

	// Calculate latency
	bestBid.Latency = time.Since(startTime).Milliseconds()

	return bestBid, nil
}

// evaluateBid evaluates a single campaign bid
func (s *AuctionService) evaluateBid(ctx context.Context, req *models.BidRequest, campaignID string) (*models.BidResponse, error) {
	// Get campaign details
	campaign, err := s.cache.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	// Check campaign status
	if campaign.Status != "active" {
		return nil, nil
	}

	// Check campaign budget
	budget, err := s.cache.GetCampaignBudget(ctx, campaignID)
	if err != nil || budget < campaign.BidPrice {
		return nil, nil
	}

	// Check targeting rules
	if !s.matchesTargeting(req, &campaign.Targeting) {
		return nil, nil
	}

	return &models.BidResponse{
		RequestID:  req.RequestID,
		CampaignID: campaignID,
		Price:      campaign.BidPrice,
		AdID:       "ad_" + campaignID,
		AdMarkup:   "<div>Sample Ad Content</div>",
	}, nil
}

// matchesTargeting checks if the request matches campaign targeting rules
func (s *AuctionService) matchesTargeting(req *models.BidRequest, targeting *models.Targeting) bool {
	// Check device type
	if len(targeting.DeviceTypes) > 0 {
		matched := false
		for _, dt := range targeting.DeviceTypes {
			if dt == req.Device.DeviceType {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Check OS
	if len(targeting.OS) > 0 {
		matched := false
		for _, os := range targeting.OS {
			if os == req.Device.OS {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Check user segments
	if len(targeting.UserSegments) > 0 {
		matched := false
		for _, segment := range targeting.UserSegments {
			for _, userSegment := range req.UserSegments {
				if segment == userSegment {
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}
