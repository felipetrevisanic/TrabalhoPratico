using HotChocolate.Types;

namespace src.graphql.inputs;

public class CreateProductInputType : InputObjectType<CreateProductInput>
{
    protected override void Configure(IInputObjectTypeDescriptor<CreateProductInput> descriptor)
    {
        descriptor.Name("CreateProductInput");
    }
}
