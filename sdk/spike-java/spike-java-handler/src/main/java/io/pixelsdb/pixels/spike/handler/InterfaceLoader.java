package io.pixelsdb.pixels.spike.handler;

import io.grpc.Server;
import io.grpc.ServerBuilder;

import java.io.File;
import java.net.URL;
import java.net.URLClassLoader;
import java.util.ServiceLoader;

public class InterfaceLoader {
    public static void main(String[] args) {
        if (args.length < 1) {
            System.err.println("请提供 JAR 文件路径作为参数");
            System.exit(1);
        }

        try {
            // 从 args 参数获取 JAR 包路径
            File jarFile = new File(args[0]);
            if (!jarFile.exists() || !jarFile.isFile()) {
                System.err.println("提供的路径不是一个有效的 JAR 文件");
                System.exit(1);
            }
//            // 打印 JAR 文件内容
//            try (JarFile jar = new JarFile(jarFile)) {
//                jar.stream().forEach(entry -> System.out.println(entry.getName()));
//            } catch (IOException e) {
//                System.err.println("无法读取 JAR 文件内容: " + e.getMessage());
//                e.printStackTrace();
//                System.exit(1);
//            }

            URL jarUrl = jarFile.toURI().toURL();
            URLClassLoader classLoader = new URLClassLoader(new URL[]{jarUrl}, InterfaceLoader.class.getClassLoader());

            // 使用 ServiceLoader 加载实现了 MyInterface 的类
            ServiceLoader<RequestHandler> serviceLoader = ServiceLoader.load(RequestHandler.class, classLoader);
            if (!serviceLoader.iterator().hasNext()) {
                System.err.println("ServiceLoader 加载失败: 没有找到实现 RequestHandler 接口的类");
                System.exit(1);
            }
            RequestHandler myInterface = serviceLoader.iterator().next();
            Server server = ServerBuilder.forPort(50052)
                    .addService(new FunctionServiceImpl(myInterface))
                    .build()
                    .start();

            System.out.println("Server started, listening on " + server.getPort());
            server.awaitTermination();
        } catch (Exception e) {
            System.err.println("ClassLoader 加载失败: " + e.getMessage());
            e.printStackTrace();
            System.exit(1);
        }
    }
}