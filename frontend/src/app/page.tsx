"use client";

import Sidebar from "../components/ChatSidebar";
import Main from "../components/ChatMainBox";
import Footer from "../components/ChatInputBox";
import { Box, Toolbar } from "@mui/material";
import AppBarNavigation from "../components/AppBarNavigation";
import ChatProvider from "../context/ChatProvider";

import {
  CssBaseline,
  ThemeProvider,
  createTheme,
  StyledEngineProvider,
} from "@mui/material";

import "../css/global.css";
import { blueGrey, grey } from "@mui/material/colors";

const muiTheme = createTheme({
  palette: {
    primary: grey,
    secondary: {
      main: blueGrey[500],
    },
    mode: "light",
  },
});

export default function Home() {
  return (
    <StyledEngineProvider injectFirst>
      <ThemeProvider theme={muiTheme}>
        <CssBaseline />
        <Box flexDirection="column" display="flex">
          <ChatProvider>
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
          </ChatProvider>
        </Box>
      </ThemeProvider>
    </StyledEngineProvider>
  );
}
