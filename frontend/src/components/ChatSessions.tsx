import React, { useContext, useEffect, useState } from "react";
import styled from "styled-components";
import {
  Avatar,
  Button,
  FormControl,
  InputLabel,
  Input,
  Container,
  Box,
} from "@mui/material";

import AddChatIcon from "@mui/icons-material/AddComment";
import { store } from "../utils/store";
import { NEW_CHAT_SESSION } from "../utils/storeConstants";
import { initChat } from "../api/chatAPI";

const ChatSessionTitle = styled.h2`
  fontsize: 18px;
  color: #4a4a4a;
  margin: 0;
`;

const ChatSessionDescription = styled.p`
  fontsize: 14px;
  color: #9b9b9b;
  margin: 0;
`;

const NewChatForm = styled.form`
  display: flex;
  flexdirection: column;
  margintop: 15px;
`;

type ChatSessionsProps = {};

const ChatSessions: React.FC<ChatSessionsProps> = () => {
  const [isNewChatOpen, setIsNewChatOpen] = useState(false);
  const [newRole, setNewRole] = useState("");
  const { state, dispatch } = useContext(store);

  const handleNewSessionSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    //Send request to backend to initialize chat session
    console.info("Requesting to start new chat session");
    const payload = { id: "", role: newRole, messages: [] };
    //wait for initChat to return before dispatching
    initChat(payload).then((response) => {
      console.info("Response from initChat: ", response);
      response.id !== "" &&
        dispatch({ type: NEW_CHAT_SESSION, payload: response });
    });

    setIsNewChatOpen(false);
    setNewRole("");
  };

  useEffect(() => {
    console.log("ChatSessions useEffect: ", state.chatSessions);
  }, [state.chatSessions]);

  return (
    <Container style={{ padding: "20px" }}>
      {state.chatSessions.map((session) => (
        <Box
          key={session.id}
          sx={{
            display: "flex",
            alignItems: "center",
            marginBottom: "15px",
            "&:hover": {
              cursor: "pointer",
              backgroundColor: "#ebebeb",
              opacity: [0.9, 0.8, 0.7],
            },
          }}
        >
          <Avatar alt="Avatar" src="https://picsum.photos/id/1011/50/50" />
          <Box>
            <ChatSessionTitle>{session.role}</ChatSessionTitle>
            <ChatSessionDescription>
              This is the last message.
            </ChatSessionDescription>
          </Box>
        </Box>
      ))}
      {/* conditional rendering of new chat form */}
      {isNewChatOpen ? (
        <NewChatForm onSubmit={handleNewSessionSubmit}>
          <Box display="flex" flexDirection="column" p="20px">
            <FormControl>
              <InputLabel htmlFor="new-chat-role">
                Enter ChatGPT role for the new chat...
              </InputLabel>
              <Input
                id="new-chat-role"
                value={newRole}
                onChange={(event) => setNewRole(event.target.value)}
              />
            </FormControl>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              size="small"
              style={{ marginTop: "15px" }}
            >
              Start
            </Button>
          </Box>
        </NewChatForm>
      ) : (
        <Box display="flex" flexDirection="column" p="20px">
          <Button
            variant="contained"
            startIcon={<AddChatIcon />}
            onClick={() => setIsNewChatOpen(true)}
          >
            New Chat
          </Button>
        </Box>
      )}
    </Container>
  );
};

export default ChatSessions;
