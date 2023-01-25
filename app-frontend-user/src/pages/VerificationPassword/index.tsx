import React from "react";
import { Box, Flex, Text, Button, HStack } from "@chakra-ui/react";
import { useForm, Controller } from "react-hook-form";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useMutation } from "react-query";

import { useShowError, useShowSuccess } from "utils/hooks";
import { getGenericErrors } from "utils/error";
import { setAuth } from "utils/function";
import { verifyAccount } from "apis";

import { Input } from "components";

type FormData = {
  password: string;
  confirm_password: string;
};

const VerificationPassword = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const state = searchParams.get("state");

  const showSuccess = useShowSuccess();
  const showError = useShowError();

  const { control, formState, handleSubmit, watch } = useForm<FormData>({
    mode: "onChange",
  });
  const { errors } = formState;

  const { mutate, isLoading } = useMutation(verifyAccount);

  const handleReset = async (values: FormData) => {
    mutate(
      { state, password: values.password },
      {
        onSuccess: (data) => {
          showSuccess("Reset Password successfully!!!");
          setAuth(data.access_token);
          navigate("/login");
        },
        onError: (error: any) => {
          showError("Reset password", getGenericErrors(error));
        },
      }
    );
  };

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
          Reset Password
        </Text>
        <Box as="form" onSubmit={handleSubmit(handleReset)}>
          <Controller
            render={({ field }) => (
              <Input
                label="New Password:"
                type="password"
                isDisabled={isLoading}
                {...field}
                errors={errors}
                containerProps={{ isRequired: true }}
              />
            )}
            control={control}
            name="password"
            rules={{
              required: "Password is required",
              minLength: {
                value: 8,
                message: "Password must have at least 8 characters",
              },
            }}
          />
          <Controller
            render={({ field }) => (
              <Input
                label="Confirm Password:"
                type="password"
                isDisabled={isLoading}
                {...field}
                errors={errors}
                containerProps={{ isRequired: true }}
              />
            )}
            control={control}
            name="confirm_password"
            rules={{
              required: "Confirm Password is required",
              validate: (val) => {
                if (watch("password") !== val)
                  return "Your passwords do no match";
              },
            }}
          />

          <HStack>
            <Button isDisabled={isLoading} variant="primary" type="submit">
              Verify
            </Button>
          </HStack>
        </Box>
      </Box>
    </Flex>
  );
};

export default VerificationPassword;
