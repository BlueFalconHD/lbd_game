"use client";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { FormEvent, useContext, useState } from "react";
import { AuthContext } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";
import { login } from "@/lib/api/auth";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircleIcon, Loader2 } from "lucide-react";
import { SetToken } from "@/lib/auth";

export function LoginForm({
  className,
  ...props
}: React.ComponentPropsWithoutRef<"form">) {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const { isAuthenticated, processLogin } = useContext(AuthContext);
  const router = useRouter();

  async function handleLogin(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    if (isAuthenticated) {
      setIsLoading(false);
      router.push("/");
      return;
    }

    try {
      const result = await login({ username: username, password: password });

      if (result.error) {
        setError(
          result.error || "An unknown error occurred. Please try again later.",
        );
        setIsLoading(false);
      } else {
        if (!result.token) {
          setError(
            "An error occurred. Please try again later or email hdombroski28@[school email domain]. Code: no_token_in_response",
          );
          setIsLoading(false);
          return;
        }

        SetToken(result.token);
        processLogin();
        setIsLoading(false);
        router.push("/");
      }
    } catch (error) {
      setError(
        "An error occurred. Please try again later or email hdombroski28@[school email domain]. Error: " +
          error,
      );
      setIsLoading(false);
    }
  }

  return (
    <AlertDialog>
      <form
        className={cn("flex flex-col gap-6", className)}
        {...props}
        onSubmit={handleLogin}
      >
        <div className="flex flex-col items-center gap-2 text-center">
          <h1 className="text-2xl font-bold">Login to your account</h1>
          <p className="text-balance text-sm text-muted-foreground">
            Enter your username below to login to your account
          </p>
        </div>
        <div>
          {error && (
            <Alert variant="destructive">
              <AlertCircleIcon className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
        </div>
        <div className="grid gap-6">
          <div className="grid gap-2">
            <Label htmlFor="username">Username</Label>
            <Input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="JohnAppleseed"
              required
            />
          </div>
          <div className="grid gap-2">
            <div className="flex items-center">
              <Label htmlFor="password">Password</Label>
              <AlertDialogTrigger asChild>
                <a
                  href="#"
                  className="ml-auto text-sm underline-offset-4 hover:underline"
                >
                  Forgot your password?
                </a>
              </AlertDialogTrigger>
            </div>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <Button type="submit" className="w-full">
            {isLoading ? (
              <>
                <Loader2 className="animate-spin" />
                Logging In...
              </>
            ) : (
              "Log in"
            )}
          </Button>
        </div>
        <div className="text-center text-sm">
          Don&apos;t have an account?{" "}
          <a href="/signup" className="underline underline-offset-4">
            Sign up
          </a>
        </div>
      </form>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Forgotten Password</AlertDialogTitle>
          <AlertDialogDescription>
            If you have forgotten your password, please contact Hayes Dombroski
            (hdombroski28@lausanneschool.com) to have it reset.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogAction>k, cool!</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
