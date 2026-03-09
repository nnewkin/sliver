fn main() {
    // 暂时注释掉 protobuf 编译，后续再添加
    // tonic_build::configure()
    //     .build_server(true)
    //     .build_client(true)
    //     .compile(&["protobuf/sliver.proto"], &["protobuf/"])
    //     .expect("Failed to compile protobuf files");
}
