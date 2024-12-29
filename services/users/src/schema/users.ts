import { relations } from "drizzle-orm";
import { AnyPgColumn, integer, pgTable, uuid, varchar } from "drizzle-orm/pg-core";
import { teams } from "./teams";

export const users = pgTable("users", {
    id: uuid("id").primaryKey(),
    name: varchar("user_name",{ length: 255}).notNull(),
    email: varchar("user_email", {length: 255}).notNull(),
    role: varchar("user_role", {length: 10}).default('admin').notNull().$type<'admin'|'user'|'moderator'>(),
    password: varchar("user_password", {length: 255}).notNull(),
    points: integer("user_points").default(0),
    teamID: uuid("team_id").references(():AnyPgColumn => teams.id,{
      onDelete: "set null"
    })
  })
export const userRelations = relations(users, ({ one,many }) => ({
  team: one(teams, { fields: [users.teamID], references: [teams.id]}),
  teamMembers: many(users)
}));
