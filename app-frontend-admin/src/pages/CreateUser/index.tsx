import React from "react";
import { useNavigate } from "react-router-dom";
import { useForm, Controller } from "react-hook-form";
import { Box, HStack, Flex, Button } from "@chakra-ui/react";
import { useMutation } from "react-query";

import { useShowError, useShowSuccess } from "utils/hooks";
import { getGenericErrors } from "utils/error";
import { createUser } from "apis";

import { Input } from "components";

type FormData = {
  name?: string;
  email: string;
  password: string;
  confirm_password: string;
  wallet: string;
};

const CreateUser = () => {
  const navigate = useNavigate();

  const showError = useShowError();
  const showSuccess = useShowSuccess();

  const { control, formState, handleSubmit, watch } = useForm<FormData>({
    mode: "onChange",
  });
  const { errors } = formState;

  const { mutate, isLoading } = useMutation(createUser);

  const handleSignUp = async (values: FormData) => {
    mutate(values, {
      onSuccess: () => {
        showSuccess("Create user Successfully!!!");
        navigate("/");
      },
      onError: (error: any) => {
        showError("Create user Failed", getGenericErrors(error));
      },
    });
  };

  return (
    <Flex
      flexDir="column"
      width={{ base: "100%", md: "70%", lg: "50%" }}
      justifyContent="center"
      alignItems="center"
      paddingBottom="5"
      margin="0 auto"
    >
      <Box
        w="100%"
        paddingX="2.5rem"
        as="form"
        onSubmit={handleSubmit(handleSignUp)}
      >
        <Controller
          render={({ field }) => (
            <Input
              label="Dislay Name:"
              isDisabled={isLoading}
              {...field}
              errors={errors}
            />
          )}
          control={control}
          name="name"
        />
        <Controller
          render={({ field }) => (
            <Input
              label="Email address:"
              isDisabled={isLoading}
              {...field}
              errors={errors}
              containerProps={{ isRequired: true }}
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
        <Controller
          render={({ field }) => (
            <Input
              label="Wallet address:"
              isDisabled={isLoading}
              {...field}
              errors={errors}
              containerProps={{ isRequired: true }}
            />
          )}
          control={control}
          name="wallet"
          rules={{
            required: "Wallet address is required",
          }}
        />
        <HStack justifyContent="center">
          <Button isDisabled={isLoading} variant="primary" type="submit">
            Create User
          </Button>
          <Button isDisabled={isLoading} onClick={() => navigate("/")}>
            Cancel
          </Button>
        </HStack>
      </Box>
    </Flex>
  );
};

export default CreateUser;
