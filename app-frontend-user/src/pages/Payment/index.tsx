import React, { useCallback } from "react";
import { Text, Box, HStack, Button } from "@chakra-ui/react";
import { useForm, Controller } from "react-hook-form";
import { useNavigate, useParams } from "react-router-dom";
import { useMutation } from "react-query";

import { updateTaskPayment } from "apis";
import { getGenericErrors } from "utils/error";
import { useShowError, useShowSuccess } from "utils/hooks";

import { Input } from "components";

type FormData = {
  transactionId: string;
};

const Payment = () => {
  const navigate = useNavigate();
  const { id } = useParams();

  const showError = useShowError();
  const showSuccess = useShowSuccess();

  const { control, formState, handleSubmit, reset } = useForm<FormData>({
    mode: "onChange",
  });
  const { errors } = formState;

  const { mutate, isLoading } = useMutation(updateTaskPayment);

  const handleSubmitNow = useCallback(
    (values: FormData) => {
      mutate(
        { taskId: id as string, paymentTxn: values.transactionId },
        {
          onSuccess: () => {
            showSuccess("Submit Payment Successfully");
            reset();
          },
          onError: (error: any) => {
            showError("Submit Payment", getGenericErrors(error));
          },
          onSettled: () => {
            reset();
            navigate("/procure");
          },
        }
      );
    },
    [id, mutate, navigate, reset, showError, showSuccess]
  );

  return (
    <div>
      <Text>Price: 1 AWO tokens</Text>
      <Text>
        Please transfer to the following wallet address:
        0x0EED72B06a3737aC6D88CEE445c6b716027b84c3
      </Text>
      <Box as="form" onSubmit={handleSubmit(handleSubmitNow)} pt="10">
        <Controller
          render={({ field }) => (
            <Input label="Transaction ID: " {...field} errors={errors} />
          )}
          control={control}
          name="transactionId"
          rules={{
            required: "This field is required",
          }}
        />

        <HStack w="full" justify="center" mt="4">
          <Button variant="primary" type="submit" isLoading={isLoading}>
            Submit Now
          </Button>
          <Button
            variant="primary"
            onClick={() => navigate("/procure")}
            isDisabled={isLoading}
          >
            Submit Later
          </Button>
        </HStack>
      </Box>
    </div>
  );
};

export default Payment;
