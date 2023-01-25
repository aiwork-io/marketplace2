import React from "react";
import {
  Box,
  MenuButton,
  Menu,
  MenuList,
  MenuItem,
  IconButton,
} from "@chakra-ui/react";
import { HamburgerIcon } from "@chakra-ui/icons";
import { logout } from "utils/function";

const Header = () => {
  return (
    <Box textAlign="right">
      <Menu>
        <MenuButton
          as={IconButton}
          aria-label="home"
          icon={<HamburgerIcon />}
          variant="outline"
          color="primary"
          rounded="full"
          colorScheme="secondary"
        />
        <MenuList>
          <MenuItem onClick={() => logout()}>Logout</MenuItem>
        </MenuList>
      </Menu>
    </Box>
  );
};

export default Header;
