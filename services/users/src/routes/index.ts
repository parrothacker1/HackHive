import { FastifyInstance } from "fastify";
import UserRouter from "./users";
import TeamRouter from "./teams";

const Router = async(fastify: FastifyInstance) => {
  fastify.register(UserRouter,{prefix:'/users'})
  fastify.register(TeamRouter,{prefix:'/teams'})
  fastify.get('/users/healthcheck',async (_) => {
    return {status:"to_be_determined"}
  })
}

export default Router
