import React, { useCallback, Suspense } from "react";
import {
  BrowserRouter as Router,
  Routes as Switch,
  Route,
} from "react-router-dom";

import { LoadingPage } from "pages";
import { Layout, PrivateRoute } from "components";

const NotFoundPage = React.lazy(() => import("pages/NotFound"));
const LoginPage = React.lazy(() => import("pages/Login"));
const HomePage = React.lazy(() => import("pages/Home"));
const CreateUserPage = React.lazy(() => import("pages/CreateUser"));

const AUTHENTICATED_ROUTES = [
  {
    path: "/",
    component: <HomePage />,
  },
  {
    path: "/create-user",
    component: <CreateUserPage />,
  },
];

const Routes = () => {
  const generateAuthenticatedRoutes = useCallback(() => {
    return AUTHENTICATED_ROUTES.map((route) => (
      <Route key={route.path} path={route.path} element={<PrivateRoute />}>
        <Route path={route.path} element={route.component} />
      </Route>
    ));
  }, []);
  return (
    <Router>
      <Suspense fallback={<LoadingPage />}>
        <Switch>
          <Route path="*" element={<NotFoundPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/" element={<Layout />}>
            {generateAuthenticatedRoutes()}
          </Route>
        </Switch>
      </Suspense>
    </Router>
  );
};

export default Routes;
