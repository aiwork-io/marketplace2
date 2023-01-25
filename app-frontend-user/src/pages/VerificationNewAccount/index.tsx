import React, { useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { useMutation } from "react-query";
import { Center, Text } from "@chakra-ui/react";

import { verifyAccount } from "apis";

import { Loading } from "components";

const VerificationNewAccount = () => {
  const [searchParams] = useSearchParams();
  const state = searchParams.get("state");

  const { mutate, isSuccess } = useMutation(verifyAccount);

  useEffect(() => {
    if (state) {
      mutate({ state });
    }
  }, [mutate, state]);

  if (isSuccess)
    return (
      <Center>
        <Text>Your email has been verified</Text>
      </Center>
    );

  return <Loading />;
};

export default VerificationNewAccount;
