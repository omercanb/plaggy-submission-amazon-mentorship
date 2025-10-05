import { sha256 } from "js-sha256";
export async function sha256Hex(input: string): Promise<string> {
  // SHA256 hash to send to the backend (hopefully)
  return sha256(input);
  // The previous implementation using Web Crypto API had issues with
  // http so I switched to js-sha256 fuck you ~brcrt
}
