import { FunctionComponent } from "react";
import Sidebar from "../components/ChatSidebar";
import AppBarNavigation from "../components/AppBarNavigation";
import Main from "../components/ChatMainBox";
import Footer from "../components/ChatInputBox";
import { Box, Toolbar } from "@mui/material";
import AppBarNavigation2 from "../components/AppBarNavigation";

const ChatPage: FunctionComponent = () => {
  return (
    <Box flexDirection="column" display="flex">
      <AppBarNavigation2 />
      <Box display="flex" height="100vh">
        <Box sx={{ flexGrow: 0, display: { xs: "none", md: "flex" } }}>
          <Sidebar />
        </Box>
        <Box width="100vh" display="flex" flexDirection="column" maxWidth="xl" flexGrow={1}>
          <Toolbar />
          <Main />
          <Footer />
        </Box>
      </Box>
    </Box>
  );
};

export default ChatPage;
