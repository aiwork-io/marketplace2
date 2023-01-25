import React, { useCallback, useEffect, useMemo, useState } from "react";
import {
  Box,
  Button,
  HStack,
  Text,
  Select,
  Input,
  Stack,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
  Center,
  ButtonProps,
} from "@chakra-ui/react";
import { useNavigate } from "react-router-dom";
import { useQuery } from "react-query";
import { isEmpty } from "lodash";
import { format, formatISO } from "date-fns";
import { useForm, Controller } from "react-hook-form";
import {
  Paginator,
  Previous,
  usePaginator,
  Next,
  PageGroup,
} from "chakra-paginator";
import { ChevronLeftIcon, ChevronRightIcon } from "@chakra-ui/icons";

import { apis, procureStatus } from "utils/config";
import { generateTaskStatus } from "utils/function";
import { getTasks } from "apis";
import { TaskStatus } from "types/task";

import { Loading } from "components";

type FormData = {
  id?: string;
  startDate?: string;
  endDate?: string;
  status?: TaskStatus;
};

const outerLimit = 2;
const innerLimit = 2;

const Procure = () => {
  const navigate = useNavigate();

  const [total, setTotal] = useState(0);
  const [filter, setFilter] = useState<{
    id?: string;
    createStart?: string;
    createEnd?: string;
    status?: TaskStatus;
  }>();

  const { control, handleSubmit, reset, watch } = useForm<FormData>({
    mode: "onChange",
  });
  const startDateWatch = watch("startDate");
  const endDateWatch = watch("endDate");

  const {
    isDisabled,
    pagesQuantity,
    currentPage,
    setCurrentPage,
    pageSize,
    offset,
  } = usePaginator({
    total: total,
    initialState: {
      pageSize: 20,
      currentPage: 1,
      isDisabled: false,
    },
  });

  const { data, isLoading } = useQuery(
    [
      apis.tasks().key,
      filter?.id,
      filter?.createStart,
      filter?.createEnd,
      filter?.status,
      pageSize,
      offset,
    ],
    () =>
      getTasks({
        ...filter,
        pageNumber: currentPage,
        pageSize: pageSize,
      })
  );

  const handleApply = (values: FormData) => {
    setFilter({
      createStart: values.startDate
        ? formatISO(new Date(values.startDate))
        : undefined,
      createEnd: values.endDate
        ? formatISO(new Date(values.endDate))
        : undefined,
      status: values.status,
    });
  };

  const handleReset = () => {
    setFilter({});
    reset();
  };

  const handlePageChange = useCallback(
    (nextPage: number) => {
      setCurrentPage(nextPage);
    },
    [setCurrentPage]
  );

  const baseStyles: ButtonProps = useMemo(
    () => ({
      w: 7,
    }),
    []
  );

  const normalStyles: ButtonProps = useMemo(
    () => ({
      ...baseStyles,
      _hover: {
        bg: "primary-light",
        color: "white",
      },
    }),
    [baseStyles]
  );

  const activeStyles: ButtonProps = useMemo(
    () => ({
      ...baseStyles,
      _hover: {
        bg: "primary-light",
        color: "white",
      },
      bg: "primary-light",
      color: "white",
    }),
    [baseStyles]
  );

  const separatorStyles: ButtonProps = useMemo(
    () => ({
      w: 7,
      bg: "primary-light",
    }),
    []
  );

  useEffect(() => {
    setTotal(data?.count || 0);
  }, [data?.count]);

  // @ts-ignore this lib is incompatible with react18
  const renderContent = useCallback(() => {
    if (isLoading) return <Loading />;
    if (isEmpty(data?.data)) return <Center>No data</Center>;
    return (
      <>
        <TableContainer>
          <Table>
            <Thead>
              <Tr>
                <Th>ID</Th>
                <Th>Name</Th>
                <Th>Submitted Date</Th>
                <Th>Status</Th>
              </Tr>
            </Thead>
            <Tbody>
              {data?.data?.map((item) => (
                <Tr
                  key={item.id}
                  onClick={() => navigate(item.id)}
                  cursor="pointer"
                  _hover={{
                    background: "primary-light",
                    color: "white",
                    transition: "0.3s all",
                  }}
                >
                  <Td>{item.id}</Td>
                  <Td>{item.name}</Td>
                  <Td>{format(new Date(item.created_at), "yyyy-MM-dd")}</Td>
                  <Td>{generateTaskStatus(item.status)}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </TableContainer>
        <Paginator
          isDisabled={isDisabled}
          activeStyles={activeStyles}
          innerLimit={innerLimit}
          currentPage={currentPage}
          outerLimit={outerLimit}
          normalStyles={normalStyles}
          separatorStyles={separatorStyles}
          pagesQuantity={pagesQuantity}
          onPageChange={handlePageChange}
        >
          <HStack pt="4" justify="flex-end">
            <Previous>
              <ChevronLeftIcon />
            </Previous>
            <PageGroup isInline align="center" />
            <Next>
              <ChevronRightIcon />
            </Next>
          </HStack>
        </Paginator>
      </>
    );
  }, [
    data,
    isLoading,
    navigate,
    activeStyles,
    currentPage,
    handlePageChange,
    isDisabled,
    normalStyles,
    pagesQuantity,
    separatorStyles,
  ]);

  return (
    <Box>
      <Button variant="primary" onClick={() => navigate("/procure/create-new")}>
        New AI compute task
      </Button>
      <Box marginTop="4" p="4" border="1px">
        <Stack
          direction={["column", "row"]}
          as="form"
          onSubmit={handleSubmit(handleApply)}
        >
          <Controller
            render={({ field }) => <Input placeholder="ID" {...field} />}
            control={control}
            name="id"
          />

          <HStack>
            <Text>Start</Text>
            <Controller
              render={({ field }) => (
                <Input max={endDateWatch} type="date" {...field} />
              )}
              control={control}
              name="startDate"
            />
          </HStack>
          <HStack>
            <Text>End</Text>
            <Controller
              render={({ field }) => (
                <Input min={startDateWatch} type="date" {...field} />
              )}
              control={control}
              name="endDate"
            />
          </HStack>
          <Controller
            render={({ field }) => (
              <Select placeholder="Status" {...field}>
                {procureStatus.map((item) => (
                  <option key={item.value} value={item.value}>
                    {item.label}
                  </option>
                ))}
              </Select>
            )}
            control={control}
            name="status"
          />

          <Button variant="primary" minW="70px" type="submit">
            Apply
          </Button>
          <Button
            variant="primary"
            minW="70px"
            type="reset"
            onClick={handleReset}
          >
            Reset
          </Button>
        </Stack>
        <Box marginTop="10">{renderContent()}</Box>
      </Box>
    </Box>
  );
};

export default Procure;
