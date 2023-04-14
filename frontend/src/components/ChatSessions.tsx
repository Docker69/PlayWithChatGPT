import React, { useContext, useState } from 'react';
import styled from 'styled-components';
import {
  Avatar,
  Button,
  FormControl,
  InputLabel,
  Input,
  Container,
  Box,
} from '@mui/material';

import AddChatIcon from '@mui/icons-material/AddComment';
import { store } from '../utils/store';


const ChatSessionTitle = styled.h2`
  fontSize: 18px;
  color: #4a4a4a;
  margin: 0;
`;

const ChatSessionDescription = styled.p`
  fontSize: 14px;
  color: #9b9b9b;
  margin: 0;
`;

const NewChatForm = styled.form`
  display: flex;
  flexDirection: column;
  marginTop: 15px;
`;

type ChatSessionsProps = {
  
};

const ChatSessions: React.FC<ChatSessionsProps> = () => {
  const [isNewChatOpen, setIsNewChatOpen] = useState(false);
  const [newChatTitle, setNewChatTitle] = useState('');
  const { dispatch } = useContext(store);

  const handleNewChatSubmit = (event: React.FormEvent<HTMLFormElement>) => {
  
    event.preventDefault();
    // handle adding new chat Session
    dispatch({
        type: "NEW_CHAT_SESSION", 
        payload: {id: newChatTitle, messages: []} 
      });
  
    console.info('new chat Session added')
    setIsNewChatOpen(false);
    setNewChatTitle('');
  };

  return (
    <Container style={{ padding: '20px' }}>
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          marginBottom: '15px',
        
          '&:hover': {
            cursor: 'pointer',
            backgroundColor: '#ebebeb',
            opacity: [0.9, 0.8, 0.7],
          },
        }}>
        <Avatar
        alt="Avatar"
        src="https://picsum.photos/id/1011/50/50"
      />
      <Box>
        <ChatSessionTitle>John Doe</ChatSessionTitle>
        <ChatSessionDescription>This is the last message.</ChatSessionDescription>
      </Box>
    </Box>

      {/* conditional rendering of new chat form */ }
  {
    isNewChatOpen ? (
      <NewChatForm onSubmit={handleNewChatSubmit}>
        <FormControl>
          <InputLabel htmlFor="new-chat-title">
            Enter a title for the new chat...
          </InputLabel>
          <Input
            id="new-chat-title"
            value={newChatTitle}
            onChange={(event) => setNewChatTitle(event.target.value)}
          />
        </FormControl>
        <Button type="submit" variant="contained" color="primary">
          Create New Chat
        </Button>
      </NewChatForm>
    ) : (
      <Button variant="contained" startIcon={<AddChatIcon />} onClick={() => setIsNewChatOpen(true)}>
        New Chat
      </Button>
    )
  }
    </Container >
  );
};

export default ChatSessions;
