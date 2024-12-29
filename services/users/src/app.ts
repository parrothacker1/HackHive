import Fastify from "fastify";
import Router from "./routes";

const fastify = Fastify({
  logger: true
});
fastify.register(Router,{prefix:'/api'})

const start = async () => {
  try {
    fastify.listen({port:3000})
  } catch(err) {
    fastify.log.error(err)
    process.exit(1)
  }
}
start()
