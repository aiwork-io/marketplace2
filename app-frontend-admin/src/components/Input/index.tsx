import React, { ReactElement } from "react";
import {
  FormControl,
  FormLabel,
  FormControlProps,
  Input as InputRaw,
  InputGroup,
  InputLeftElement,
  InputRightElement,
  InputProps,
  Text,
  Textarea,
} from "@chakra-ui/react";

interface MyInputProps extends InputProps {
  label?: string;
  name?: string;
  errors?: any;
  leftElement?: ReactElement;
  rightElement?: ReactElement;
  containerProps?: FormControlProps;
  componentType?: "text" | "textarea";
}

const Input = React.forwardRef<HTMLInputElement, MyInputProps>(
  (
    {
      label,
      componentType = "text",
      errors,
      name,
      leftElement,
      rightElement,
      containerProps = {},
      isRequired,
      ...others
    }: MyInputProps,
    ref
  ) => {
    const error = name && errors?.[name]?.message;
    const Component = componentType === "text" ? InputRaw : Textarea;

    return (
      <FormControl mb="1rem" {...containerProps} isInvalid={!!error}>
        {!!label && (
          <FormLabel fontSize="1.125rem" fontWeight="400" mb="0.5rem">
            {label}
          </FormLabel>
        )}
        <InputGroup>
          {!!leftElement && (
            <InputLeftElement pointerEvents="none">
              {leftElement}
            </InputLeftElement>
          )}
          {/* @ts-ignore */}
          <Component name={name} ref={ref} {...others} isInvalid={!!error} />
          {!!rightElement && (
            <InputRightElement>{rightElement}</InputRightElement>
          )}
        </InputGroup>
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
  }
);

export default Input;
