"use client";

import * as React from "react";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import {
  NEW_CHAT_SESSION,
  SET_HUMAN,
  SET_NEW_CHAT_DIALOG_STATE,
} from "../global/ChatProviderConstants";
import { ChatContext } from "../context/ChatProvider";
import { initChat, initSession } from "../app/api/chatAPI";

export default function NewChatDialog() {
  const { state, dispatch } = React.useContext(ChatContext);
  const [newRole, setNewRole] = React.useState("");

  const handleNewSessionSubmit = () => {
    //event.preventDefault();
    //Send request to backend to initialize chat session
    console.info("Requesting to start new chat session");
    const payload = {
      id: "",
      role: newRole,
      humanId: state.human.id,
      messages: [],
    };
    //wait for initChat to return before dispatching
    initChat(payload).then(({ success, response }) => {
      console.info("Response from initChat: ", { success, response });
      if (success) {
        dispatch({ type: NEW_CHAT_SESSION, payload: response });

        //get update chat list
        initSession(state.human).then(({ success, response }) => {
          console.info("Response from initSession: ", { success, response });
          //TODO: handle the case where nickname not found to create new human
          success &&
            response.id !== "" &&
            dispatch({ type: SET_HUMAN, payload: response });
        });
      }
    });

    handleCancel();
    setNewRole("");
  };

  const handleCancel = () => {
    dispatch({
      type: SET_NEW_CHAT_DIALOG_STATE,
      payload: false,
    });
  };

  return (
    <div>
      <Dialog open={state.chatDialogOpen} onClose={handleCancel}>
        <DialogTitle>Create New Chat Session</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Enter the desired role for the new chat. For example, &quot;Customer
            Service&quot; or &quot;Helpful Developer&quot;
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="name"
            label="ChatGPT Role"
            type="text"
            fullWidth
            variant="standard"
            onChange={(v) => setNewRole(v.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCancel}>Cancel</Button>
          <Button onClick={handleNewSessionSubmit}>Create</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
