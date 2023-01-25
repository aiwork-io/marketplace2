import React, { useCallback, useMemo, useState } from "react";
import { useForm, Controller } from "react-hook-form";
import {
  Box,
  Button,
  Center,
  Select,
  FormLabel,
  FormControl,
  HStack,
  Text,
} from "@chakra-ui/react";
import { useNavigate } from "react-router-dom";
import { chunk, every } from "lodash";
import { useMutation } from "react-query";

import { createTask, uploadImage, uploadImageToStore } from "apis";
import { delay } from "utils/function";
import { useShowError, useShowSuccess } from "utils/hooks";
import { getGenericErrors } from "utils/error";

import { Input, InputFile } from "components";
import { categoryOption } from "utils/config";

type FormData = {
  name: string;
  category: string;
  files: FileList;
};

const CHUNK = 5;

const CreateTask = () => {
  const navigate = useNavigate();

  const { control, formState, handleSubmit, reset } = useForm<FormData>({
    mode: "onChange",
  });
  const { errors } = formState;

  const [isUploading, setIsUploading] = useState(false);

  const showError = useShowError();
  const showSuccess = useShowSuccess();

  const { mutate: handleCreateTask, isLoading: isCreatingTask } =
    useMutation(createTask);

  const handleSave = useCallback(
    (values: FormData) => {
      setIsUploading(true);
      handleCreateTask(
        { name: values.name, category: values.category },
        {
          onSuccess: (data) => {
            const taskId = data.id;
            const imagesChunk = chunk(Array.from(values.files), CHUNK);
            Promise.all(
              imagesChunk.map(async (images: File[], index) => {
                await delay(index * 500).then(async () => {
                  await Promise.all(
                    Array.from(images).map(async (image) => {
                      try {
                        const res = await uploadImage({
                          taskId,
                          filename: image.name,
                        });
                        await uploadImageToStore({ file: image, url: res.url });
                      } catch (error: any) {
                        return showError(
                          "Create Task",
                          getGenericErrors(error)
                        );
                      }
                    })
                  );
                });
              })
            ).then(() => {
              setIsUploading(true);
              delay(3000).then(() => {
                setIsUploading(false);
                reset();
                showSuccess("Created successfully!");
                navigate(`/procure/${taskId}/payment`);
              });
            });
          },
        }
      );
    },
    [handleCreateTask, navigate, reset, showError, showSuccess]
  );

  const isLoading = useMemo(
    () => isCreatingTask || isUploading,
    [isCreatingTask, isUploading]
  );

  return (
    <Box as="form" onSubmit={handleSubmit(handleSave)}>
      <Controller
        render={({ field }) => (
          <Input label="AI compute task name: " {...field} errors={errors} />
        )}
        control={control}
        name="name"
        rules={{
          required: "This field is required",
        }}
      />
      <Controller
        render={({ field }) => {
          const error = errors?.["name"]?.message;
          return (
            <FormControl mb="1rem">
              <FormLabel fontSize="1.125rem" fontWeight="400" mb="0.5rem">
                Category:
              </FormLabel>
              <Select
                {...field}
                placeholder="Select Category"
                isInvalid={!!error}
              >
                {categoryOption.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </Select>
              {!!error && (
                <Text
                  whiteSpace="pre-wrap"
                  mt="0.5rem"
                  fontSize="0.75rem"
                  fontWeight="400"
                  color="error"
                >
                  {error}
                </Text>
              )}
            </FormControl>
          );
        }}
        control={control}
        name="category"
        rules={{
          required: "This field is required",
        }}
      />
      <Controller
        control={control}
        name="files"
        render={({ field }) => (
          <FormControl mb="1rem">
            <FormLabel fontSize="1.125rem" fontWeight="400" mb="0.5rem">
              File(s):
            </FormLabel>
            <InputFile accept=".jpg,.png" {...field} errors={errors} multiple>
              <Button variant="primary" w="full">
                Browse
              </Button>
            </InputFile>
            {field?.value?.length && (
              <Center>Files: {field?.value?.length}</Center>
            )}
          </FormControl>
        )}
        rules={{
          required: "This field is required",
          validate: (val) => {
            if (every(val, (item) => item.size > 1 * 1024 * 1024)) {
              return "Max file size is 80MB";
            }
          },
        }}
      />
      <Text>
        * Accepted file types: jpg, png. Max file size: 80MB. Upload images that
        you would like to be processed
      </Text>
      <HStack w="full" justify="center" mt="4">
        <Button variant="primary" type="submit" isLoading={isLoading}>
          Submit
        </Button>
        <Button onClick={() => navigate("/procure")} isDisabled={isLoading}>
          Cancel
        </Button>
      </HStack>
    </Box>
  );
};

export default CreateTask;
