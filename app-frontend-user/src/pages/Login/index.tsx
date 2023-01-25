import React from "react";
import { useForm, Controller } from "react-hook-form";
import { Link, useNavigate, Navigate } from "react-router-dom";
import { Box, Flex, Text, Button, HStack } from "@chakra-ui/react";
import { useMutation } from "react-query";

import { useShowError, useShowSuccess } from "utils/hooks";
import { isLoggedIn, setAuth } from "utils/function";
import { getGenericErrors } from "utils/error";
import { login } from "apis";

import { Input } from "components";

type FormData = {
  email: string;
  password: string;
};

const Login = () => {
  const navigate = useNavigate();

  const showError = useShowError();
  const showSuccess = useShowSuccess();
  const { control, formState, handleSubmit } = useForm<FormData>({
    mode: "onChange",
  });
  const { errors } = formState;

  const loggedIn = isLoggedIn();

  const { mutate, isLoading } = useMutation(login);

  const handleLogin = async (values: FormData) => {
    mutate(values, {
      onSuccess: (data) => {
        // if (!data.verified_at) {
        //   return showError("Login Failed", "your email have not been verify");
        // }
        setAuth(data.access_token);
        showSuccess("Login successfully!");
      },
      onError: (error: any) => {
        showError("Login Failed", getGenericErrors(error));
      },
    });
  };

  if (loggedIn) return <Navigate to="/" />;

  return (
    <Flex
      flexDir="column"
      width={{ base: "100%", md: "70%", lg: "50%" }}
      height="100vh"
      justifyContent="center"
      alignItems="center"
      paddingBottom="5"
      margin="0 auto"
    >
      <Box w="100%" paddingX="2.5rem">
        <Text
          fontWeight={700}
          fontSize="4xl"
          textAlign="center"
          marginBottom="10"
        >
          Marketplace
        </Text>
        <Box as="form" onSubmit={handleSubmit(handleLogin)}>
          <Controller
            render={({ field }) => (
              <Input
                label="Email address:"
                isDisabled={isLoading}
                {...field}
                errors={errors}
              />
            )}
            control={control}
            name="email"
            rules={{
              required: "Email is required",
              pattern: {
                value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                message: "Invalid email address",
              },
            }}
          />
          <Controller
            render={({ field }) => (
              <Input
                label="Password:"
                type="password"
                isDisabled={isLoading}
                {...field}
                errors={errors}
              />
            )}
            control={control}
            name="password"
            rules={{
              required: "Password is required",
            }}
          />
          <HStack justifyContent="center">
            <Button isDisabled={isLoading} variant="primary" type="submit">
              Login
            </Button>
            <Button
              isDisabled={isLoading}
              variant="primary"
              onClick={() => navigate("/sign-up")}
            >
              Register
            </Button>
          </HStack>
          <Box textAlign="center" paddingTop="5">
            <Text
              as={Link}
              to="/forgot-password"
              textDecor="underline"
              fontWeight="bold"
            >
              Forget Password
            </Text>
          </Box>
        </Box>
      </Box>
    </Flex>
  );
};

export default Login;
