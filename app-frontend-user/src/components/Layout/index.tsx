import React from "react";
import { Box, Container } from "@chakra-ui/react";
import { Outlet } from "react-router-dom";

import { Header } from "components";

const Layout = () => {
  return (
    <Container maxW="1000px" paddingY="10">
      <Header />
      <Box paddingY="10">
        <Outlet />
      </Box>
    </Container>
  );
};

export default Layout;
