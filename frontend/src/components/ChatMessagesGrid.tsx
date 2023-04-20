import { Grid, Avatar, Typography, styled } from "@mui/material";

interface ChatClasses {
  [key: string]: string;
  avatar: string;
  msg: string;
  leftRow: string;
  left: string;
  leftFirst: string;
  leftLast: string;
  rightRow: string;
  right: string;
  rightFirst: string;
  rightLast: string;
}

interface ChatMessagesProps {
  avatar: string;
  messages: string[];
  side: "left" | "right";
}

const classes: ChatClasses = {
  avatar: `ChatMessages-avatar`,
  msg: `ChatMessages-msg`,
  leftRow: `ChatMessages-leftRow`,
  left: `ChatMessages-left`,
  leftFirst: `ChatMessages-leftFirst`,
  leftLast: `ChatMessages-leftLast`,
  rightRow: `ChatMessages-rightRow`,
  right: `ChatMessages-right`,
  rightFirst: `ChatMessages-rightFirst`,
  rightLast: `ChatMessages-rightLast`,
};

const StyledGrid = styled(Grid)(({ theme: { palette, spacing } }) => {
  const radius = spacing(2.5);
  const size = spacing(4);
  const rightBgColor = palette.primary.main;
  // if you want the same as facebook messenger, use this color '#09f'
  return {
    [`& .${classes.avatar}`]: {
      width: size,
      height: size,
    },
    [`& .${classes.msg}`]: {
      padding: spacing(1, 2),
      borderRadius: 4,
      marginBottom: 4,
      display: "inline-block",
      wordBreak: "break-all",
      fontFamily:
        // eslint-disable-next-line max-len
        '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol"',
    },
    [`& .${classes.leftRow}`]: {
      textAlign: "left",
    },
    [`& .${classes.left}`]: {
      borderTopRightRadius: radius,
      borderBottomRightRadius: radius,
      backgroundColor: palette.grey[100],
    },
    [`& .${classes.leftFirst}`]: {
      borderTopLeftRadius: radius,
    },
    [`& .${classes.leftLast}`]: {
      borderBottomLeftRadius: radius,
    },
    [`& .${classes.rightRow}`]: {
      textAlign: "right",
    },
    [`& .${classes.right}`]: {
      borderTopLeftRadius: radius,
      borderBottomLeftRadius: radius,
      backgroundColor: rightBgColor,
      color: palette.common.white,
    },
    [`& .${classes.rightFirst}`]: {
      borderTopRightRadius: radius,
    },
    [`& .${classes.rightLast}`]: {
      borderBottomRightRadius: radius,
    },
  };
});

const ChatMessagesGrid = ({ avatar, messages, side }: ChatMessagesProps) => {
  const attachClass = (index: number): string => {
    if (index === 0) {
      return classes[`${side}First`];
    }
    if (index === messages.length - 1) {
      return classes[`${side}Last`];
    }
    return "";
  };

  return (
    <StyledGrid
      container
      spacing={2}
      justifyItems={side === "right" ? "flex-end" : "flex-start"}
    >
      {side === "left" && (
        <Grid item>
          <Avatar className={classes.avatar} src={avatar} />
        </Grid>
      )}
      <Grid item xs={11}>
        {messages.map((msg, i) => (
          <div key={i} className={classes[`${side}Row`]}>
            <Typography
              align={"left"}
              className={`${classes.msg} ${classes[`${side}`]} ${attachClass(
                i
              )}`}
            >
              {msg}
            </Typography>
          </div>
        ))}
      </Grid>
    </StyledGrid>
  );
};

export default ChatMessagesGrid;
