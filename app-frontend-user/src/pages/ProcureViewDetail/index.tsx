import React, { useCallback, useMemo } from "react";
import {
  Box,
  Button,
  SimpleGrid,
  VStack,
  Text,
  Center,
} from "@chakra-ui/react";
import { useNavigate, useParams } from "react-router-dom";
import { useMutation, useQuery } from "react-query";
import { isEmpty } from "lodash";
import { format } from "date-fns";

import { downloadResult, getTaskDetail } from "apis";
import { apis } from "utils/config";
import { generateTaskStatus } from "utils/function";
import { useShowError, useShowSuccess } from "utils/hooks";
import { getGenericErrors } from "utils/error";
import { TaskStatus } from "types/task";

import { Loading } from "components";
import ResultImage from "./ResultImage";

const ProcureViewDetail = () => {
  const navigate = useNavigate();
  const { id } = useParams();

  const showSuccess = useShowSuccess();
  const showError = useShowError();

  const { data, isLoading } = useQuery(
    [apis.task().key, id],
    () => getTaskDetail(id as string),
    {
      enabled: !!id,
    }
  );

  const { mutate, isLoading: isDownloading } = useMutation(downloadResult);

  const results = useMemo(() => {
    return data?.assets.map((item) => JSON.parse(item.results || "{}"));
  }, [data?.assets]);

  const handleDownloadResult = useCallback(() => {
    mutate(id as string, {
      onSuccess: (data) => {
        const link = document.createElement("a");
        link.href = data.url;
        link.download = "result";
        link.click();
        showSuccess("Download successfully!!!");
      },
      onError: (error: any) => {
        showError("Download", getGenericErrors(error));
      },
    });
  }, [id, mutate, showError, showSuccess]);

  const render = useCallback(() => {
    if (isLoading) return <Loading />;
    if (isEmpty(data)) return <Center>No data</Center>;
    return (
      <>
        <SimpleGrid columns={[1, 3]} spacing="4">
          <Box>ID: {data.id}</Box>
          <Box>Name: {data.name}</Box>
          <Box>
            Submitted Date: {format(new Date(data.created_at), "yyyy-MM-dd")}
          </Box>
          <Box>Status: {generateTaskStatus(data.status)}</Box>
          <Box>
            <VStack align="flex-start">
              <Button
                isDisabled={data.status !== TaskStatus.PENDING}
                onClick={() => navigate(`/procure/${id}/payment`)}
              >
                Submit Payment Proof
              </Button>
              <Button
                isDisabled={data.status !== TaskStatus.COMPLETED}
                onClick={() => handleDownloadResult()}
                isLoading={isDownloading}
              >
                Download json result
              </Button>
            </VStack>
          </Box>
        </SimpleGrid>
        <Box marginTop="4">
          <Text>Result</Text>
          {results?.map((item, index) => (
            <ResultImage key={index} data={item} />
          ))}
        </Box>
      </>
    );
  }, [
    data,
    handleDownloadResult,
    id,
    isDownloading,
    isLoading,
    navigate,
    results,
  ]);

  return (
    <Box>
      <Button variant="primary" onClick={() => navigate("/procure/create-new")}>
        New AI compute task
      </Button>
      <Box marginTop="4" p="4" border="1px">
        {render()}
      </Box>
    </Box>
  );
};

export default ProcureViewDetail;
