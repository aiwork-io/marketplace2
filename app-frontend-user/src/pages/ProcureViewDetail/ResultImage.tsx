import React, { useEffect, useRef } from "react";
import { Box, Center } from "@chakra-ui/react";
import { isEmpty, round } from "lodash";

import { Result } from "types/task";

type ResultImageProps = {
  data: Result;
};

const ResultImage = ({ data }: ResultImageProps) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const boxRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const ctx = canvasRef.current?.getContext("2d");
    const image = new Image();
    image.onload = () => {
      if (ctx) {
        const boxWidth = boxRef.current?.offsetWidth || 0;
        const canvas = ctx?.canvas;
        const widthImage = image.width;
        const heightImage = image.height;
        const ratio = boxWidth / widthImage;

        const width = widthImage * ratio;
        const height = heightImage * ratio;

        if (canvas) {
          canvas.width = width;
          canvas.height = height;
        }

        ctx.drawImage(image, 0, 0, width, height);

        const colors = ["red", "green", "blue", "orange", "gray"];

        for (let i = 0; i < data.data.object.length; i++) {
          const color = colors[i % colors.length];
          const item = data.data.object[i];
          const text = `${item.category} ${round(item.score, 2)}`;
          const box = item.bbox;
          ctx.beginPath();
          ctx.rect(
            box[0] * ratio,
            box[1] * ratio,
            box[2] * ratio,
            box[3] * ratio
          );
          ctx.lineWidth = 4;
          ctx.strokeStyle = color;
          ctx.closePath();

          let length = text.length;
          length = length * 14 * 0.62;
          ctx.fillStyle = color;
          ctx.fillRect(
            box[0] * ratio - 2,
            box[1] * ratio - 14 * 1.5,
            length,
            14 * 1.5
          );
          ctx.font = "14px monospace";
          ctx.fillStyle = "white";

          ctx.fillText(text, box[0] * ratio, box[1] * ratio - 4);

          ctx.stroke();
        }
      }
    };
    image.src = data.source;
  }, [data?.data?.object, data?.source]);

  if (isEmpty(data)) return <Center>No data</Center>;

  return (
    <Box ref={boxRef}>
      <canvas ref={canvasRef} />
    </Box>
  );
};

export default ResultImage;
