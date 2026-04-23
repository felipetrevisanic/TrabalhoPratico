using HotChocolate.Types;

namespace src.interfaces.graphql.inputs;

public class CreateProductInputType : InputObjectType<CreateProductInput>
{
    protected override void Configure(IInputObjectTypeDescriptor<CreateProductInput> descriptor)
    {
        descriptor.Name("CreateProductInput");
    }
}
