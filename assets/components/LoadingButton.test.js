import LoadingButton from '@/components/LoadingButton';
import {render} from '@testing-library/vue';

describe('LoadingButton', () => {
  it('renders the text from the default slot', async () => {
    const {findByText} = render(LoadingButton, {
      slots: {
        default: 'Hello there',
      },
    });

    await findByText('Hello there');
  });

  it('doesn\'t show the default slot while in loading state', async () => {
    const {queryByText} = render(LoadingButton, {
      propsData: {
        loading: true,
      },
      slots: {
        default: 'Hello there',
      },
    });

    expect(await queryByText('Hello there')).toBeNull();
  });

  it('renders a custom loading state', async () => {
    const {findByText} = render(LoadingButton, {
      propsData: {
        loading: true,
      },
      slots: {
        loading: 'I am loading!',
      },
    });

    await findByText('I am loading!');
  });

  it('button is not disabled if prop disabled is falsy', async () => {
    const {findByRole} = render(LoadingButton, {
      propsData: {
        disabled: false,
      },
    });

    expect(await findByRole('button')).not.toBeDisabled();
  });

  it('disables the button if disabled prop ist set', async () => {
    const {findByRole} = render(LoadingButton, {
      propsData: {
        disabled: true,
      },
    });

    expect(await findByRole('button')).toBeDisabled();
  });

  it('sets the correct variant class', async () => {
    const {findByRole} = render(LoadingButton, {
      propsData: {
        variant: 'danger',
      },
    });

    expect(await findByRole('button')).toHaveClass('btn-danger');
  });
});
