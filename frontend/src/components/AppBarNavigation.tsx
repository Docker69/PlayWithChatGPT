"use client"

import { FunctionComponent, useContext, useState } from "react";
import {
  AppBar,
  Avatar,
  Box,
  IconButton,
  Menu,
  MenuItem,
  Toolbar,
  Tooltip,
  Typography,
} from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import React from "react";
import { ChatContext } from "../context/ChatProvider";
import { SET_DRAWER_STATE } from "../global/ChatProviderConstants";
import { signOut } from "next-auth/react";

interface UserMenuProps {
  settings: string[];
  anchorEl: HTMLElement | null;
  onClose?: React.MouseEventHandler | undefined;

}

const UserMenu: FunctionComponent<UserMenuProps> = ({
  settings,
  anchorEl,
  onClose,
}) => (
  <Menu
    sx={{ mt: "45px" }}
    id="menu-appbar"
    anchorEl={anchorEl}
    anchorOrigin={{
      vertical: "top",
      horizontal: "right",
    }}
    keepMounted
    transformOrigin={{
      vertical: "top",
      horizontal: "right",
    }}
    open={Boolean(anchorEl)}
    onClose={onClose}
  >
    {settings.map((setting) => (
      <MenuItem key={setting} onClick={onClose}>
        <Typography textAlign="center">{setting}</Typography>
      </MenuItem>
    ))}
  </Menu>
);

const settings = ["Profile", "Account", "Dashboard", "Logout"];

const AppBarNavigation2: FunctionComponent = () => {
  const [anchorElUser, setAnchorElUser] = useState<null | HTMLElement>(null);

  const handleOpenUserMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElUser(event.currentTarget);
  };
  
  const handleCloseUserMenu = (event: React.MouseEvent<HTMLElement>) => {
    if (event.currentTarget.textContent==='Logout')
    {
      signOut();
    }
    setAnchorElUser(null);
  };

  const { state, dispatch } = useContext(ChatContext);

  // Set the state of Component B
  const setDrawerState = () => {
    dispatch({ type: SET_DRAWER_STATE, payload: !state.mobileDrawerOpen });
  };

  return (
    <AppBar
      position="fixed"
      sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}
    >
      <div>
        <Toolbar disableGutters>
          <Box
            sx={{
              ml: 2,
              display: { xs: "none", md: "flex" },
              mr: 1,
              justifySelf: "flex-start",
            }}
          >
            <img alt="" src="/logo50.svg" />
          </Box>
          <Typography
            variant="h6"
            noWrap
            component="a"
            href="/"
            sx={{
              mr: 2,
              display: { xs: "none", md: "flex" },
              fontFamily: "monospace",
              fontWeight: 700,
              letterSpacing: ".3rem",
              color: "inherit",
              textDecoration: "none",
            }}
          >
            Play With ChatGPT
          </Typography>

          <Box sx={{ flexGrow: 1, display: { xs: "flex", md: "none" } }}>
            <IconButton
              size="large"
              aria-label="account of current user"
              aria-controls="menu-appbar"
              aria-haspopup="true"
              onClick={setDrawerState}
              color="inherit"
            >
              <MenuIcon />
            </IconButton>
          </Box>
          <Box sx={{ display: { xs: "flex", md: "none" }, mr: 1 }}></Box>

          <Typography
            variant="h6"
            noWrap
            component="a"
            href=""
            sx={{
              mr: 2,
              display: { xs: "flex", md: "none" },
              flexGrow: 1,
              fontFamily: "monospace",
              fontWeight: 700,
              letterSpacing: ".3rem",
              color: "inherit",
              textDecoration: "none",
            }}
          >
            Play With ChatGPT
          </Typography>
          <Box sx={{ flexGrow: 1, display: { xs: "none", md: "flex" } }}></Box>
          <Box sx={{ flexGrow: 0 }} marginRight={2}>
            <Tooltip title="Open settings">
              <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                <Avatar alt="Remy Sharp" src="/static/images/avatar/2.jpg" />
              </IconButton>
            </Tooltip>
            <UserMenu
              settings={settings}
              anchorEl={anchorElUser}
              onClose={handleCloseUserMenu}
            />
          </Box>
        </Toolbar>
      </div>
    </AppBar>
  );
};

export default AppBarNavigation2;
