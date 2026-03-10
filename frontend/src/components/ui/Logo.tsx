/**
 * Logo component that switches between color and white versions
 * based on the user's color scheme preference.
 */

import Image from "next/image";

interface LogoProps {
  width?: number;
  height?: number;
  className?: string;
}

export function Logo({ width = 120, height = 40, className = "" }: LogoProps) {
  return (
    <>
      <Image
        src="/logo.svg"
        alt="Loupi"
        width={width}
        height={height}
        className={`dark:hidden ${className}`}
        priority
      />
      <Image
        src="/logo-white.svg"
        alt="Loupi"
        width={width}
        height={height}
        className={`hidden dark:block ${className}`}
        priority
      />
    </>
  );
}
