import { FunctionComponent } from "react";
import Sidebar from "../components/ChatSidebar";
import Main from "../components/ChatMainBox";
import Footer from "../components/ChatInputBox";
import { Box, Toolbar } from "@mui/material";
import AppBarNavigation from "../components/AppBarNavigation";
import ChatProvider from "../context/ChatProvider";

//TODO: Memoize the provider and children to prevent unnecessary re-renders
const ChatPage: FunctionComponent = () => {

  return (
    <ChatProvider>
      <Box flexDirection="column" display="flex">
        <AppBarNavigation />
        <Box display="flex" height="100vh">
          <Box sx={{ flexGrow: 0, display: { xs: "none", md: "flex" } }}>
            <Sidebar />
          </Box>
          <Box
            width="100vh"
            height="100vh"
            display="flex"
            flexDirection="column"
            maxWidth="xl"
            flexGrow={1}
          >
            <Toolbar />
            <Main />
            <Footer />
          </Box>
        </Box>
      </Box>
    </ChatProvider>
  );
};

export default ChatPage;
