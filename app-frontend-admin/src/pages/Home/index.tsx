import React from "react";
import { Center, Button } from "@chakra-ui/react";
import { useNavigate } from "react-router-dom";

const Home = () => {
  const navigate = useNavigate();
  return (
    <Center>
      <Button onClick={() => navigate("/create-user")} variant="primary">
        Create user account
      </Button>
    </Center>
  );
};

export default Home;
