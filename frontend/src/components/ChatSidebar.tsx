"use client";

import {
  FunctionComponent,
  memo,
  useContext,
  useEffect,
  useState,
} from "react";
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
  Human,
  SEND_CHAT_PROMPT,
  SET_DRAWER_STATE,
  SET_INIT_SESSION_DIALOG_STATE,
  SET_NEW_CHAT_DIALOG_STATE,
} from "../global/ChatProviderConstants";
import NewChatDialog from "../dialogs/NewChatDialog";
import { ChatContext } from "../context/ChatProvider";
import IdentifyDialog from "../dialogs/IdentifyDialog";
import React from "react";
import { getChatSession } from "@/app/api/chatAPI";
import { stat } from "fs";

const ChatSidebar: FunctionComponent = () => {
  const drawerWidth = 240;
  const { state, dispatch } = useContext(ChatContext);

  console.debug("ChatSidebar render");

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

  const handleSetChatSession = (value: string | Human) => {
    //wait for initSession to return before dispatching
    const id = typeof value === "string" ? value : value.chatIds[0].id;
    getChatSession(id).then(({ success, response }) => {
      console.info("Response from getChatSession: ", { success, response });
      //TODO: handle the case where nickname not found to create new human
      success &&
        response.id !== "" &&
        dispatch({ type: SEND_CHAT_PROMPT, payload: response });
    });
  };

  useEffect(() => {
    console.log(
      "ChatSidebar useEffect, state.mobileDrawerOpen",
      state.mobileDrawerOpen
    );
  }, [state.mobileDrawerOpen]);

  useEffect(() => {
    console.debug("ChatSidebar useEffect");
    state.human.id === "" &&
      dispatch({ type: SET_INIT_SESSION_DIALOG_STATE, payload: true });

    state.human.id !== "" && handleSetChatSession(state.human);
  }, [state.human.id]);

  const drawer = (
    <div>
      <Toolbar />
      <Divider />
      <List dense={true} disablePadding>
        {state.human.chatIds.map((chat) => (
          //if the session is active, set the active class
          <ListItemButton
            key={chat.id}
            selected={chat.id === state.activeChatSession.id}
            //prevent click if waiting for response
            disabled={state.waitingForResponse}
            onClick={() => handleSetChatSession(chat.id)}
          >
            <ListItemIcon sx={{ color: "inherit" }}>
              <ChatIcon />
            </ListItemIcon>
            <ListItemText
              primary={chat.role}
              primaryTypographyProps={{
                noWrap: true,
                fontSize: 14,
                fontWeight: "medium",
              }}
            />
          </ListItemButton>
        ))}
        <ListItemButton
          key="New Chat"
          onClick={openNewChatDialog}
          disabled={state.waitingForResponse}
        >
          <ListItemIcon>
            <AddCommentIcon />
          </ListItemIcon>
          <ListItemText primary={"New Chat"} />
        </ListItemButton>
      </List>
      <NewChatDialog />
      <IdentifyDialog />
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
