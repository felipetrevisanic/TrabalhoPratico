using HotChocolate.Types;

namespace src.graphql.inputs;

public class UpdateProductInputType : InputObjectType<UpdateProductInput>
{
    protected override void Configure(IInputObjectTypeDescriptor<UpdateProductInput> descriptor)
    {
        descriptor.Name("UpdateProductInput");
    }
}
