import { AnyPgColumn, integer, jsonb, pgTable, uuid, varchar } from "drizzle-orm/pg-core";
import { users } from "./users";
import { relations } from "drizzle-orm";

export const teams = pgTable("teams",{
    id: uuid("id").primaryKey(),
    name: varchar("team_name",{length: 255}).notNull(),
    members:jsonb().$type<string[]>().notNull(),
    points: integer("points").default(0).notNull(),
    leader: uuid("team_leader").references(():AnyPgColumn => users.id,{
      onDelete:"set null"
    }).notNull()
})
export const teamRelations = relations(teams,({ one }) => ({
  leader: one(users, {
    fields: [teams.leader],
    references: [users.id]
  })
}))  
