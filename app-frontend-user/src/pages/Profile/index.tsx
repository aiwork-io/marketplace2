import React, { useEffect } from "react";
import { Box, HStack, Button, Text } from "@chakra-ui/react";
import { useForm, Controller } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { useMutation, useQuery } from "react-query";

import { useShowError, useShowSuccess } from "utils/hooks";
import { getGenericErrors } from "utils/error";

import { Input } from "components";
import { getProfile, updateProfile } from "apis";
import { apis } from "utils/config";

type FormData = {
  name: string;
  wallet_address: string;
};

const Profile = () => {
  const navigate = useNavigate();
  const showSuccess = useShowSuccess();
  const showError = useShowError();

  const { control, formState, handleSubmit, setValue } = useForm<FormData>({
    mode: "onChange",
  });
  const { errors } = formState;

  const { data, isLoading } = useQuery(apis.getProfile.key, getProfile);

  const { mutate, isLoading: isUpdating } = useMutation(updateProfile);

  const handleSave = async (values: FormData) => {
    console.log("handleSave");
    mutate(
      { name: values.name, wallet: values.wallet_address },
      {
        onSuccess: () => {
          showSuccess("Updated profile successfully!");
        },
        onError: (error: any) => {
          showError("Updated profile failed!", getGenericErrors(error));
        },
      }
    );
  };

  useEffect(() => {
    setValue("name", data?.name || "");
    setValue("wallet_address", data?.wallet || "");
  }, [data?.name, data?.wallet, setValue]);

  const loading = isLoading || isUpdating;

  return (
    <Box w="100%" as="form" onSubmit={handleSubmit(handleSave)}>
      <Controller
        render={({ field }) => (
          <Input label="Display Name:" {...field} errors={errors} />
        )}
        control={control}
        name="name"
        rules={{
          required: "Display name is required",
        }}
      />
      <HStack mb="4">
        <Text fontSize="1.125rem" fontWeight="400" mb="0.5rem">
          Password:{" "}
        </Text>
        <Button
          isDisabled={loading}
          variant="primary"
          onClick={() => navigate("/forgot-password")}
        >
          Reset password
        </Button>
      </HStack>
      <Controller
        render={({ field }) => (
          <Input
            label="Wallet address:"
            isDisabled={loading}
            {...field}
            errors={errors}
          />
        )}
        control={control}
        name="wallet_address"
        rules={{
          required: "Wallet address is required",
        }}
      />
      <HStack justifyContent="center" marginTop="4">
        <Button isDisabled={loading} variant="primary" type="submit">
          Save
        </Button>
        <Button isDisabled={loading} onClick={() => navigate(-1)}>
          Back
        </Button>
      </HStack>
    </Box>
  );
};

export default Profile;
