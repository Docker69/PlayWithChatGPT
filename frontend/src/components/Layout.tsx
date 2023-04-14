import { Box } from '@mui/material';
import ChatSessions from './ChatSessions';

type LayoutProps = {
  children: React.ReactNode;
};

const Layout: React.FC<LayoutProps> = ({ children }) => {  return (
    <Box display="flex" flexDirection="row">
      <Box width="25%" height="100vh">
        <ChatSessions />
      </Box>
      <Box width="75%" height="100vh">
        {children}
      </Box>
    </Box>
  );
};

export default Layout;
