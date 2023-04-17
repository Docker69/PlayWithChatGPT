import { Box } from '@mui/material';
import ChatSessions from './ChatSessions';
import ChatSession from './ChatSession';
import { FunctionComponent } from 'react';

const Layout: FunctionComponent = () => {  return (
    <Box display="flex" flexDirection="row" width={"inherit"}>
      <Box width="25%" height="100vh">
        <ChatSessions />
      </Box>
      <Box width="75%" height="100vh">
        <ChatSession />
      </Box>
    </Box>
  );
};

export default Layout;
