import React from "react";
import { IconButton, Box, Container, TextField } from "@mui/material";
import SendIcon from "@mui/icons-material/Send";

type ChatSessionProps = {};

const MessageBubble = ({
  author,
  message,
  mine,
}: {
  author: string;
  message: string;
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
      {author}
    </p>
    <p style={{ fontSize: "14px", margin: "0" }}>{message}</p>
  </Box>
);

const ChatSession: React.FC<ChatSessionProps> = () => {
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

  return (
    <Container
      maxWidth="sm"
      style={{ height: "100%", display: "flex", flexDirection: "column" }}
    >
      <Box sx={{ flex: "1", overflowY: "auto", p: "20px" }}>
        {messages.map((message, index) => (
          <MessageBubble key={index} {...message} />
        ))}
      </Box>
      <Box display="flex" p="20px">
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Type your message..."
          style={{ flex: 1, marginRight: "10px" }}
        />
        <IconButton color="primary" aria-label="Send" size="large">
          <SendIcon />
        </IconButton>
      </Box>
    </Container>
  );
};

export default ChatSession;
