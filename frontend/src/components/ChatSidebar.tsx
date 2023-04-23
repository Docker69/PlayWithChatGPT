"use client"

import { FunctionComponent, useContext, useEffect } from "react";
import {
  Divider,
  Drawer,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
} from "@mui/material";
import AddCommentIcon from "@mui/icons-material/AddComment";
import ChatIcon from "@mui/icons-material/Chat";
import {
  SET_DRAWER_STATE,
  SET_NEW_CHAT_DIALOG_STATE,
} from "../global/ChatProviderConstants";
import NewChatDialog from "../dialogs/NewChatDialog";
import { ChatContext } from "../context/ChatProvider";

const ChatSidebar: FunctionComponent = () => {
  console.debug("ChatSidebar render");

  const drawerWidth = 240;
  const { state, dispatch } = useContext(ChatContext);

  // Set the state of the Mobile Drawer
  const setDrawerState = () => {
    dispatch({ type: SET_DRAWER_STATE, payload: !state.mobileDrawerOpen });
  };

  // Set the state of the New Chat Dialog
  const openNewChatDialog = () => {
    dispatch({
      type: SET_NEW_CHAT_DIALOG_STATE,
      payload: true,
    });
  };

  useEffect(() => {
    console.log(
      "ChatSidebar useEffect, state.mobileDrawerOpen",
      state.mobileDrawerOpen
    );
  }, [state.mobileDrawerOpen]);

  useEffect(() => {
    console.log(
      "ChatSidebar useEffect, state.chatSessions: ",
      state.chatSessions
    );
  }, [state.chatSessions]);

  const drawer = (
    <div>
      <Toolbar />
      <Divider />
      <List>
        {state.chatSessions.map((session) => (
          //if the session is active, set the active class
          <ListItemButton
            key={session.id}
            selected={session.id === state.activeChatSession.id}
          >
            <ListItemIcon>
              <ChatIcon />
            </ListItemIcon>
            <ListItemText primary={session.role} />
          </ListItemButton>
        ))}
        <ListItemButton key="New Chat" onClick={openNewChatDialog}>
          <ListItemIcon>
            <AddCommentIcon />
          </ListItemIcon>
          <ListItemText primary={"New Chat"} />
        </ListItemButton>
      </List>
      <NewChatDialog />
    </div>
  );

  return (
    <>
      <Drawer
        variant="temporary"
        open={state.mobileDrawerOpen}
        onClose={setDrawerState}
        ModalProps={{
          keepMounted: true, // Better open performance on mobile.
        }}
        sx={{
          display: { xs: "block", md: "none" },
          "& .MuiDrawer-paper": { boxSizing: "border-box", width: drawerWidth },
        }}
      >
        {drawer}
      </Drawer>
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: {
            width: drawerWidth,
            boxSizing: "border-box",
          },
        }}
      >
        {drawer}
      </Drawer>
    </>
  );
};

export default ChatSidebar;
