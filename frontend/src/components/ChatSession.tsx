import React, { useContext, useEffect, useState } from "react";
import { IconButton, Box, Container, TextField } from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import { store } from "../utils/store";
import { SEND_CHAT_PROMPT, USER_ROLE } from "../utils/storeConstants";
import { sendChatPrompt } from "../api/chatAPI";
type ChatSessionProps = {};

const MessageBubble = ({
  role,
  content,
  mine,
}: {
  role: string;
  content: string;
  mine: boolean;
}) => (
  <Box
    maxWidth="60%"
    bgcolor={mine ? "#ebebeb" : "#f4f4f4"}
    alignSelf={mine ? "flex-end" : "flex-start"}
    borderRadius="15px"
    padding="8px 12px"
    marginBottom="15px"
  >
    <p
      style={{
        fontSize: "12px",
        color: "#9b9b9b",
        margin: "0 0 5px",
        textAlign: mine ? "right" : "left",
      }}
    >
      {role}
    </p>
    <p style={{ fontSize: "14px", margin: "0" }}>{content}</p>
  </Box>
);

const ChatSession: React.FC<ChatSessionProps> = () => {
/*
  const messages = [
    { author: "John Doe", message: "Hello!", mine: false },
    { author: "Jane Smith", message: "Hi, how are you?", mine: true },
    {
      author: "John Doe",
      message: "I am doing good. Thanks for asking.",
      mine: false,
    },
    {
      author: "Jane Smith",
      message: "How is your day going?",
      mine: true,
    },
    {
      author: "John Doe",
      message: "It is going well. How about yours?",
      mine: false,
    },
  ];
*/
  const [prompt, setPrompt] = useState("");
  const { state, dispatch } = useContext(store);

  const handleSendClick = () => {
    //event.preventDefault();
    //use chat api to send prompt to backend
    console.info("Requesting to send prompt to backend");

    const updatedMessages = [...state.activeChatSession.messages, {role: USER_ROLE, content: prompt}];
    const payload = {...state.activeChatSession, messages: updatedMessages };

    //wait for sendChatPrompt to return before dispatching
    sendChatPrompt(payload).then(({success, response}) => {
      console.info("Response from sendChatPrompt: ", {success, response});
      success &&
        dispatch({ type: SEND_CHAT_PROMPT, payload: response });
    });

    setPrompt("");
  };

  //use effect to update state when active chat session changes or new message is sent
  useEffect(() => {
    console.log("ChatSession useEffect: ", state.activeChatSession.messages);
  }, [state.activeChatSession.messages]);

  return (
    <Container
      maxWidth="sm"
      style={{ height: "100%", display: "flex", flexDirection: "column" }}
    >
      <Box sx={{ flex: "1", overflowY: "auto", p: "20px" }}>
        {state.activeChatSession.messages.map((message, index) => (
          <MessageBubble key={index} {...message} mine={message.role === USER_ROLE ? true : false} />
        ))}
      </Box>
      <Box display="flex" p="20px">
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Type your message..."
          style={{ flex: 1, marginRight: "10px" }}
          value={prompt}
          onChange={(event) => setPrompt(event.target.value)}
        />
        <IconButton
          color="primary"
          aria-label="Send"
          size="large"
          //handle click event in TSX and prevent default
          onClick={handleSendClick}
        >
          <SendIcon />
        </IconButton>
      </Box>
    </Container>
  );
};

export default ChatSession;
