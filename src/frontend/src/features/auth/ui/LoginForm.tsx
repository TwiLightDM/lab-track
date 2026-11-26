import { useState } from "react";
import type { ChangeEvent, FormEvent } from "react";
import { useNavigate } from "react-router-dom";

import type { ILoginData } from "@/features/auth";
import { paths } from "@/app/router/paths";

import "./AuthForms.css";

interface LoginFormProps {
  onSubmit: (data: ILoginData) => Promise<void>;
}

const LoginForm = ({ onSubmit }: LoginFormProps) => {
  const [formData, setFormData] = useState<ILoginData>({
    email: "",
    password: "",
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const [fieldErrors, setFieldErrors] = useState<{
    email?: string;
    password?: string;
  }>({});

  const navigate = useNavigate();

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));

    setError("");

    setFieldErrors((prev) => ({
      ...prev,
      [name]: undefined,
    }));
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setLoading(true);
    setError("");
    setFieldErrors({});

    try {
      await onSubmit(formData);
      navigate(paths.profile.url());
    } catch (err) {
      const message = err instanceof Error ? err.message : "Ошибка входа";

      setError(message);

      if (message.toLowerCase().includes("email")) {
        setFieldErrors({ email: message });
      } else if (
        message.toLowerCase().includes("пароль") ||
        message.toLowerCase().includes("password")
      ) {
        setFieldErrors({ password: message });
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="AuthForms__Wrapper">
      <form className="AuthForms" onSubmit={handleSubmit}>
        <h2 className="AuthForms__Title">Добро пожаловать!</h2>

        <p className="AuthForms__Subtitle">Войдите в свой аккаунт</p>

        {error && !fieldErrors.email && !fieldErrors.password && (
          <div className="AuthForms__Error">{error}</div>
        )}

        <div className="AuthForms__Group">
          <label htmlFor="email">Email</label>

          <input
            id="email"
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            placeholder="your@email.com"
            disabled={loading}
            required
            className={fieldErrors.email ? "AuthForms__Input--error" : ""}
          />

          {fieldErrors.email && (
            <div className="AuthForms__FieldError">{fieldErrors.email}</div>
          )}
        </div>

        <div className="AuthForms__Group">
          <label htmlFor="password">Пароль</label>

          <input
            id="password"
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            placeholder="********"
            disabled={loading}
            required
            className={fieldErrors.password ? "AuthForms__Input--error" : ""}
          />

          {fieldErrors.password && (
            <div className="AuthForms__FieldError">{fieldErrors.password}</div>
          )}
        </div>

        <button type="submit" className="AuthForms__Button" disabled={loading}>
          {loading ? "Вход..." : "Войти"}
        </button>
      </form>
    </div>
  );
};

export default LoginForm;
