"use client";
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';

import { signIn } from "next-auth/react";

function Login() {

  return (
    <>
      <div >
        <Box
          sx={{
            margin: 15,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
          }}
        >
          <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
            <LockOutlinedIcon />
          </Avatar>
          <Button
            onClick={() => signIn()}

            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
          >
            Sign In to use ChatGPT
          </Button>
        </Box>
      </div>
    </>



  );
}

export default Login;
