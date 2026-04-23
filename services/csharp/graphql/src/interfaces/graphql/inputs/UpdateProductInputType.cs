using HotChocolate.Types;

namespace src.interfaces.graphql.inputs;

public class UpdateProductInputType : InputObjectType<UpdateProductInput>
{
    protected override void Configure(IInputObjectTypeDescriptor<UpdateProductInput> descriptor)
    {
        descriptor.Name("UpdateProductInput");
    }
}
