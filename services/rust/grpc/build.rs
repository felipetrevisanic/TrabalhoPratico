fn main() {
    let protoc = protoc_bin_vendored::protoc_bin_path().expect("failed to fetch vendored protoc");
    unsafe {
        std::env::set_var("PROTOC", protoc);
    }

    tonic_build::configure()
        .build_server(true)
        .build_client(true)
        .compile_protos(
            &["../../../shared/proto/product/v1/product.proto"],
            &["../../../shared/proto"],
        )
        .expect("failed to compile product proto");
}
