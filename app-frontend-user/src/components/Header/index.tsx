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
import { useNavigate } from "react-router-dom";
import { logout } from "utils/function";

const Header = () => {
  const navigate = useNavigate();

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
          <MenuItem onClick={() => navigate("/profile")}>Profile</MenuItem>
          <MenuItem onClick={() => navigate("/profile")}>
            Download Client
          </MenuItem>
          <MenuItem onClick={() => navigate("/procure")}>Procure View</MenuItem>
          <MenuItem onClick={() => navigate("/contribute")}>
            Contribute View
          </MenuItem>
          <MenuItem onClick={() => logout()}>Logout</MenuItem>
        </MenuList>
      </Menu>
    </Box>
  );
};

export default Header;
