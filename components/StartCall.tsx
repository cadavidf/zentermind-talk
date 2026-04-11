import { useVoice } from "@humeai/voice-react";
import { AnimatePresence, motion } from "motion/react";
import { Button } from "./ui/button";
import { Phone } from "lucide-react";
import { toast } from "sonner";

function getErrorMessage(error: unknown) {
  if (error instanceof Error && error.message) {
    return error.message;
  }

  if (typeof error === "string" && error.trim()) {
    return error;
  }

  return "Unable to start call";
}

export default function StartCall({
  configId,
  accessToken,
  isSafariIOS,
}: {
  configId?: string;
  accessToken: string;
  isSafariIOS: boolean;
}) {
  const { status, connect } = useVoice();

  return (
    <AnimatePresence>
      {status.value !== "connected" ? (
        <motion.div
          className={"fixed inset-0 p-4 flex items-center justify-center bg-background"}
          initial="initial"
          animate="enter"
          exit="exit"
          variants={{
            initial: { opacity: 0 },
            enter: { opacity: 1 },
            exit: { opacity: 0 },
          }}
        >
          <AnimatePresence>
            <motion.div
              variants={{
                initial: { scale: 0.5 },
                enter: { scale: 1 },
                exit: { scale: 0.5 },
              }}
              className="flex flex-col items-center gap-4"
            >
              <Button
                className={"z-50 flex items-center gap-1.5 rounded-full"}
                onClick={() => {
                  connect({ 
                    auth: { type: "accessToken", value: accessToken },
                    configId, 
                    // additional options can be added here
                    // like resumedChatGroupId and sessionSettings
                  })
                    .then(() => {})
                    .catch((error: unknown) => {
                      toast.error(getErrorMessage(error));
                    })
                    .finally(() => {});
                }}
              >
                <span>
                  <Phone
                    className={"size-4 opacity-50 fill-current"}
                    strokeWidth={0}
                  />
                </span>
                <span>Start Call</span>
              </Button>
              <div className="text-gray-500 dark:text-gray-400 text-sm font-light">
                serenity beyond sanctuary
              </div>
            </motion.div>
          </AnimatePresence>
          {isSafariIOS ? (
            <div className="fixed bottom-24 left-1/2 -translate-x-1/2 w-[calc(100%-2rem)] max-w-md rounded-xl border border-border bg-card px-4 py-3 text-sm text-muted-foreground shadow-lg">
              iPhone Safari can fail on the voice audio pipeline. This build
              disables AudioWorklet on Safari/iPhone, but if the call still
              does not start, use Chrome or another browser.
            </div>
          ) : null}
        </motion.div>
      ) : null}
    </AnimatePresence>
  );
}
