export const routes = [
  {
    path: "/",
    component: HomePage,
  },
  {
    path: "/movies",
    component: MoviesPage,
  },
  {
    path: /\/movies\/(\d+)/,
    component: MovieDetailsPage,
  },
  {
    path: "/account/register",
    component: RegisterPage,
  },
  {
    path: "/account/login",
    component: LoginPage,
  },
  {
    path: "/account/",
    component: AccountPage,
    loggedIn: true,
  },
  {
    path: "/account/favorites",
    component: FavoritePage,
    loggedIn: true,
  },
  {
    path: "/account/watchlist",
    component: WatchlistPage,
    loggedIn: true,
  },
];
