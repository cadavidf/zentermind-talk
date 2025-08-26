import React, { SVGAttributes } from "react";

export type ZenterMindLogoProps = SVGAttributes<SVGSVGElement>;

export default function ZenterMindLogo(props: ZenterMindLogoProps) {
  return (
    <svg
      width="200"
      height="60"
      viewBox="0 0 200 60"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <text
        x="10"
        y="35"
        fontSize="24"
        fontWeight="bold"
        fill="currentColor"
        fontFamily="system-ui, -apple-system, sans-serif"
      >
        ZenterMind
      </text>
      <circle
        cx="180"
        cy="30"
        r="15"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
      />
      <circle cx="180" cy="30" r="8" fill="currentColor" opacity="0.6" />
      <circle cx="180" cy="30" r="4" fill="currentColor" />
    </svg>
  );
}