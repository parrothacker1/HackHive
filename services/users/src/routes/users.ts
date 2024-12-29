import { FastifyInstance } from "fastify";

const UserRouter = async(fastify: FastifyInstance) => {
  fastify.post('/',async(_) => {
    return "to make a user"
  })
  fastify.put('/',async(_) => {
    return "to update a user"
  })
  fastify.delete('/',async(_) => {
    return "to delete a user"
  })
  fastify.get('/me',async(_) => {
    return "to send details of the user"
  })
}

export default UserRouter
