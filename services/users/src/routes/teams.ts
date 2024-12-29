import { FastifyInstance } from "fastify";

const TeamRouter = async(fastify: FastifyInstance) => {
  fastify.post("/",async(_) => {
    return "to_create a new team"
  })
  fastify.put("/",async(_) => {
    return "to_upate the team data"
  })
  fastify.delete("/",async(_) => {
    return "to_delete the team"
  })
  fastify.get("/me",async(_) => {
    return "to_get details about the team"
  })
}

export default TeamRouter
