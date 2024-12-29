import { defineConfig } from "drizzle-kit";

export default defineConfig({
  dialect: "postgresql",
  driver: "pglite",
  schema: "./src/schema/*",
  out: "./migrations",
  dbCredentials: {
    url: "still_need_work"
  },
  introspect:{
    casing:"camel"
  },
  migrations:{
    schema:"public"
  }
})
