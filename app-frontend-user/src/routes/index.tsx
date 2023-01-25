import React, { useCallback, Suspense } from "react";
import {
  BrowserRouter as Router,
  Routes as Switch,
  Route,
  Navigate,
} from "react-router-dom";

import { LoadingPage } from "pages";
import { Layout, PrivateRoute } from "components";

const NotFoundPage = React.lazy(() => import("pages/NotFound"));
const LoginPage = React.lazy(() => import("pages/Login"));
const SignUpPage = React.lazy(() => import("pages/SignUp"));
const ProcurePage = React.lazy(() => import("pages/Procure"));
const ProfilePage = React.lazy(() => import("pages/Profile"));
const ForgotPasswordPage = React.lazy(() => import("pages/ForgotPassword"));
const ProcureViewDetailPage = React.lazy(
  () => import("pages/ProcureViewDetail")
);
const CreateTaskPage = React.lazy(() => import("pages/CreateTask"));
const PaymentPage = React.lazy(() => import("pages/Payment"));
const ContributePage = React.lazy(() => import("pages/Contribute"));
const VerificationPasswordPage = React.lazy(
  () => import("pages/VerificationPassword")
);
const VerificationNewAccountPage = React.lazy(
  () => import("pages/VerificationNewAccount")
);

const AUTHENTICATED_ROUTES = [
  {
    path: "/procure",
    component: <ProcurePage />,
  },
  {
    path: "/procure/create-new",
    component: <CreateTaskPage />,
  },
  {
    path: "/procure/:id/payment",
    component: <PaymentPage />,
  },
  {
    path: "/procure/:id",
    component: <ProcureViewDetailPage />,
  },
  {
    path: "/profile",
    component: <ProfilePage />,
  },
  {
    path: "/contribute",
    component: <ContributePage />,
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
          <Route path="/sign-up" element={<SignUpPage />} />
          <Route path="/forgot-password" element={<ForgotPasswordPage />} />
          <Route
            path="/verification-reset-password"
            element={<VerificationPasswordPage />}
          />
          <Route
            path="/verification-new-account"
            element={<VerificationNewAccountPage />}
          />
          <Route path="/" element={<Layout />}>
            <Route index element={<Navigate to="/procure" replace />} />
            {generateAuthenticatedRoutes()}
          </Route>
        </Switch>
      </Suspense>
    </Router>
  );
};

export default Routes;
