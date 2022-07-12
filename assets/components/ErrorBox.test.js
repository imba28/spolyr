import ErrorBox from '@/components/ErrorBox';
import {render} from '@testing-library/vue';

describe('ErrorBox', () => {
  it('renders a message provided by the message prop', () => {
    const message = 'Something went wrong';
    const {getByText} = render(ErrorBox, {
      propsData: {
        message,
      },
    });

    getByText(message);
  });

  it('Renders an optional text if provided by the text prop', () => {
    const helpText = 'An helpful text';
    const {getByText} = render(ErrorBox, {
      propsData: {
        message: 'Test',
        text: helpText,
      },
    });

    getByText(helpText);
  });

  it('Contains an image', () => {
    const {container} = render(ErrorBox, {
      propsData: {
        message: 'Test',
      },
    });

    const image = container.querySelector('img');
    expect(image).toBeVisible();
  });
});
